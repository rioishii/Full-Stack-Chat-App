"use strict";

const mongoose = require("mongoose");
const express = require("express");
const morgan = require("morgan");
const { 
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
} = require('./handlers');
const { channelSchema, messageSchema } = require('./schemas');
const amqp = require('amqplib/callback_api');
const mongoEndpoint = "mongodb://mongodb:27017/test";
const rabbitEndPoint = "amqp://rabbitmq:5672";

const app = express();

let rabbitChannel;

const getRabbitChannel = () => {
    return rabbitChannel;
}

const addr = process.env.ADDR || ":5001";	
const [host, port] = addr.split(":");

app.use(express.json());
app.use(morgan("dev"));
app.use("/", (req, res, next) => {
    const user = req.header('X-User')
    if (!user) {
        res.status(401).send("User is not authenticated");
        return;
    }
    next();
});

const Channel = mongoose.model("Channel", channelSchema);
const Message = mongoose.model("Message", messageSchema);

const connect = () => {
    mongoose.connect(mongoEndpoint);  
}

const RequestWrapper = (handler, SchemeAndDBForwarder) => {
    return (req, res) => {
        handler(req, res, SchemeAndDBForwarder);
    }
}

app.get("/v1/channels", RequestWrapper(getChannelHandler, { Channel }));
app.post("/v1/channels", RequestWrapper(postChannelHandler, { Channel, getRabbitChannel }));
app.get("/v1/channels/:channelID", RequestWrapper(getSpecificChannelHandler, { Channel, Message }));
app.post("/v1/channels/:channelID", RequestWrapper(postSpecificChannelHandler, { Channel, Message, getRabbitChannel }));
app.patch("/v1/channels/:channelID", RequestWrapper(patchSpecificChannelHandler, { Channel, getRabbitChannel }));
app.delete("/v1/channels/:channelID", RequestWrapper(deleteSpecificChannelHandler, { Channel, Message, getRabbitChannel }));
app.post("/v1/channels/:channelID/members", RequestWrapper(postNewMemberHandler, { Channel }));
app.delete("/v1/channels/:channelID/members", RequestWrapper(deleteMemberHandler, { Channel }));
app.patch("/v1/messages/:messageID", RequestWrapper(patchSpecificMessageHandler, { Channel, Message, getRabbitChannel }));
app.delete("/v1/messages/:messageID", RequestWrapper(deleteSpecificMessageHandler, { Channel, Message, getRabbitChannel }));

connect();
mongoose.connection.on('error', console.error)
    .on('disconnected', connect)
    .once('open', main)

initDB();

async function main() {
    amqp.connect(rabbitEndPoint, (err, conn) => {
        if (err) {
            console.log("Error connecting to rabbit instance");
            process.exit(1);
        }
        conn.createChannel((err, ch) => {
            if (err) {
                console.log("Error connecting to channel");
                process.exit(1);
            }
            ch.assertQueue("messageQueue", {durable: true});
            rabbitChannel = ch;
            app.listen(port, host, () => {
                console.log(`message server is running on port ${port}`);
            });
        });
    });
}

function initDB() {
    const name = "general";
    const description = "general channel";
    const isPrivate = false;
    let members = [];
    const createdAt = new Date();
    const channel = {
        name,
        description,
        isPrivate,
        members,
        createdAt
    }

    const query = new Channel(channel);
    query.save()
        .then(newChannel => {
            console.log(newChannel);
        }) 
        .catch(err => {
            console.log(err);
        });
}