const getChannelHandler = async (req, res, { Channel }) => {
    const user = JSON.parse(req.header('X-User'));
    res.set("Content-Type", "application/json");
    try {
        const channels = await Channel.find({ members: { $elemMatch: { userId: user.id } } });
        res.status(201).json(channels);
    } catch {
        res.status(500).send('Unable to get channels'); 
    }
};

const postChannelHandler = (req, res, { Channel, getRabbitChannel }) => {
    const user = JSON.parse(req.header('X-User'));

    res.set("Content-Type", "application/json");

    const { name, description, isPrivate } = req.body;
    const creator = {
        userId: user.id,
        userName: user.userName,
        firstName: user.firstName,
        lastName: user.lastName
    }
    let members = [];
    members.push(creator);
    const createdAt = new Date();
    const channel = {
        name,
        description,
        isPrivate,
        members,
        createdAt,
        creator
    }

    const query = new Channel(channel);
    query.save((err, newChannel) => {
        if (err) {
            res.status(500).send('Unable to create a channel');
        }

        let ch = getRabbitChannel();
        let userIDs;
        if (newChannel.isPrivate) {
            userIDs = newChannel.members.map(member => {
                return member.userId;
            });
        } else {
            userIDs = null;
        }
        ch.sendToQueue("messageQueue", Buffer.from(JSON.stringify(
            {
                type: "channel-new",
                channel: newChannel,
                userIDs: userIDs
            }
        )));
        res.status(201).json(newChannel);
    });
};

const getSpecificChannelHandler = async (req, res, { Channel, Message }) => {
    const user = JSON.parse(req.header('X-User'));

    res.set("Content-Type", "application/json");

    const beforeID = req.query.before;

    const channelID = req.params.channelID

    try {
        const channel = await Channel.findById(channelID)
        if (channel.isPrivate) {
            let exists = channel.members.find(member => member.userId === user.id);
            if (!exists) {
                res.status(403).send("User not part of the private channel");
                return;
            }
        }
        const messages = await Message.find({channelID: channelID});
        messages.sort((a, b) => b.createdAt - a.createdAt);
        if (beforeID) {
            let index = messages.findIndex(message => message.id == beforeID);
            if (index === -1) {
                index = 0;
            }
            let messagesSlice = []
            for (let i = index; i < index + 100 && i < messages.length; i++) {
                messagesSlice.push(messages[i]);
            }
            res.status(201).json(messagesSlice);
        } else {
            let messagesSlice = []
            for (let i = 0; i < 100 && i < messages.length; i++) {
                messagesSlice.push(messages[i]);
            }
            res.status(201).json(messagesSlice);
        }
    } catch (err) {
        console.log(err);
        res.status(500).send(err);
    }
};

const postSpecificChannelHandler = async (req, res, { Channel, Message, getRabbitChannel }) => {
    const user = JSON.parse(req.header('X-User'));

    res.set("Content-Type", "application/json");

    const channelID = req.params.channelID;

    try {
        const channel = await Channel.findById(req.params.channelID)
        if (channel.isPrivate) {
            let exists = channel.members.find(member => member.userId === user.id);
            if (!exists) {
                res.status(403).send("User not part of the private channel");
                return;
            }
        }

        const { body } = req.body;
        const createdAt = new Date();
        const creator = {
            id: user.id,
            userName: user.userName,
            firstName: user.firstName,
            lastName: user.lastName
        }
        const message = {
            channelID,
            body,
            createdAt,
            creator
        };
        const query = new Message(message);
        query.save((err, newMessage) => {
            if (err) {
                res.status(500).send('Unable to create a message');
                return;
            }

            let ch = getRabbitChannel();
            let userIDs;
            if (channel.isPrivate) {
                userIDs = channel.members.map(member => {
                    return member.userId;
                });
            } else {
                userIDs = null;
            }
            ch.sendToQueue("messageQueue", Buffer.from(JSON.stringify(
                {
                    type: "message-new",
                    message: newMessage,
                    userIDs: userIDs
                }
            )));

            res.status(201).json(newMessage);
        });
    } catch (err) {
        res.status(500).send(err);
    }
};

const patchSpecificChannelHandler = async (req, res, { Channel, getRabbitChannel }) => {
    const user = JSON.parse(req.header('X-User'));

    res.set("Content-Type", "application/json");

    const channelID = req.params.channelID;

    try {
        const channel = await Channel.findById(channelID);
        if (channel.creator.userId !== user.id) {
            res.status(403).send("User is not the creator of the channel");
            return;
        }

        const { name, description } = req.body;
        const updateObj = {name: name, description: description, editedAt: new Date()};

        Channel.findByIdAndUpdate(channelID, updateObj, {multi: true, new: true}, (err, updatedChannel) => {
            if (err) {
                res.status(500).send("There was an issue");
                return;
            }

            let ch = getRabbitChannel();
            let userIDs;
            if (updatedChannel.isPrivate) {
                userIDs = updatedChannel.members.map(member => {
                    return member.userId;
                });
            } else {
                userIDs = null;
            }
            ch.sendToQueue("messageQueue", Buffer.from(JSON.stringify(
                {
                    type: "channel-update",
                    channel: updatedChannel,
                    userIDs: userIDs
                }
            )));
            res.status(201).json(updatedChannel);
        }); 
    } catch (err) {
        res.status(500).send(err);
    }
};

const deleteSpecificChannelHandler = async (req, res, { Channel, Message, getRabbitChannel }) => {
    const user = JSON.parse(req.header('X-User'));

    res.set("Content-Type", "text/html");

    const channelID = req.params.channelID

    try {
        const channel = await Channel.findById(channelID);
        if (channel.creator.userId !== user.id) {
            res.status(403).send("User is not the creator of the channel");
            return;
        }

        Channel.remove({_id: channelID}, (err, deletedChannel) => {
            if (err) {
                res.status(500).send("Failed deleting channel");
            };
            Message.remove({channelID: channelID}, (err, _) => {
                if (err) {
                    res.status(500).send("Failed deleting messages");
                }

                let ch = getRabbitChannel();
                let userIDs;
                if (deletedChannel.isPrivate) {
                    userIDs = deletedChannel.members.map(member => {
                        return member.userId;
                    });
                } else {
                    userIDs = null;
                }
                ch.sendToQueue("messageQueue", Buffer.from(JSON.stringify(
                    {
                        type: "channel-delete",
                        channel: deletedChannel,
                        userIDs: userIDs
                    }
                )));
                res.status(200).send("channel deleted successfully");
            });
        });
    } catch (err) {
        res.status(500).send(err);
    }
};

const postNewMemberHandler = async (req, res, { Channel }) => {
    const user = JSON.parse(req.header('X-User'));

    res.set("Content-Type", "text/html");

    const channelID = req.params.channelID;

    const { id, userName, firstName, lastName } = req.body;
    const newMember = {
        id: id,
        userName: userName,
        firstName: firstName,
        lastName: lastName
    }

    try {
        const channel = await Channel.findById(channelID);
        if (channel.creator.userId !== user.id) {
            res.status(403).send("User is not the creator of the channel");
            return;
        }

        Channel.findOneAndUpdate({_id: channelID}, {$push: {members: newMember}})
            .then(() => {
                res.status(201).send("New member has been successfully added to channel");
            })
            .catch(() => {
                res.status(500).send("Unable to add new member to channel");
            });
    } catch (err) {
        res.status(500).send(err);
    }
};

const deleteMemberHandler = async (req, res, { Channel }) => {
    const user = JSON.parse(req.header('X-User'));

    res.set("Content-Type", "text/html");

    const channelID = req.params.channelID;

    const { userName } = req.body;

    try {
        const channel = await Channel.findById(channelID);
        if (channel.creator.userId !== user.id) {
            res.status(403).send("User is not the creator of the channel");
            return;
        }
        Channel.findOneAndUpdate({_id: channelID}, {$pull: { "members": { "userName": userName}}}, { safe: true })
            .then(() => {
                res.status(200).send("Member has been successfully removed from channel");
            })
            .catch(() => {
                res.status(500).send("Unable to remove member from channel");
            });
    } catch (err) {
        res.status(500).send(err);
    }
};

const patchSpecificMessageHandler = async (req, res, { Channel, Message, getRabbitChannel }) => {
    const user = JSON.parse(req.header('X-User'));

    res.set("Content-Type", "application/json");

    const messageID = req.params.messageID;

    const { body } = req.body;

    try {
        const message = await Message.findById(messageID);
        if (message.creator.userId !== user.id) {
            res.status(403).send("User is not the creator of the message");
            return;
        }
        const channel = await Channel.find({ channelID: message.channelID });
        const updateObj = {body: body, editedAt: new Date()};
        Message.findByIdAndUpdate(messageID, updateObj, {multi: true, new: true}, (err, updatedMessage) => {
            if (err) {
                res.status(500).send("There was an issue updating message");
                return;
            }
            let ch = getRabbitChannel();
            let userIDs;
            if (channel.isPrivate) {
                userIDs = channel.members.map(member => {
                    return member.userId;
                });
            } else {
                userIDs = null;
            }
            ch.sendToQueue("messageQueue", Buffer.from(JSON.stringify(
                {
                    type: "message-update",
                    message: updatedMessage,
                    userIDs: userIDs
                }
            )));
            res.status(201).json(updatedMessage);
        });
    } catch (err) {
        res.status(500).send(err);
    }

    const updateObj = {body: body, editedAt: new Date()};
    Message.findByIdAndUpdate(messageID, updateObj, {multi: true, new: true}, (err, updatedMessage) => {
        if (err) {
            res.status(500).send("There was an issue updating message");
            return;
        }
        res.status(201).json(updatedMessage);
    });
}

const deleteSpecificMessageHandler = async (req, res, { Message, getRabbitChannel }) => {
    const user = JSON.parse(req.header('X-User'));

    res.set("Content-Type", "application/json");

    const messageID = req.params.messageID;

    try {
        const message = await Message.findById(messageID)
        if (message.creator.userId !== user.id) {
            res.status(403).send("User is not the creator of the message");
            return;
        }
        const channel = await Channel.find({ channelID: message.channelID });
        Message.remove({_id: messageID}, (err, deletedMessage) => {
            if (err) {
                res.status(500).send("Issue deleting message");
            }
            let ch = getRabbitChannel();
            let userIDs;
            if (channel.isPrivate) {
                userIDs = channel.members.map(member => {
                    return member.userId;
                });
            } else {
                userIDs = null;
            }
            ch.sendToQueue("messageQueue", Buffer.from(JSON.stringify(
                {
                    type: "message-delete",
                    message: deletedMessage,
                    userIDs: userIDs
                }
            )));
            res.status(200).send("Message successfully deleted");
        });
    } catch (err) {
        res.status(500).send(err);
    }
}

module.exports = {
    getChannelHandler, 
    postChannelHandler, 
    getSpecificChannelHandler, 
    postSpecificChannelHandler, 
    patchSpecificChannelHandler,
    deleteSpecificChannelHandler,
    postNewMemberHandler,
    deleteMemberHandler,
    patchSpecificMessageHandler,
    deleteSpecificMessageHandler
};