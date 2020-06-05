const Schema = require('mongoose').Schema;

const channelSchema = new Schema({
    name: {type: String, required: true, unique: true},
    description: String,
    isPrivate: Boolean,
    members: [{
        userId: Number,
        email: String,
        userName: String,
        firstName: String,
        lastName: String
    }],
    createdAt: Date,
    creator: {
        userId: Number,
        email: String,
        userName: String,
        firstName: String,
        lastName: String
    },
    editedAt: Date
});

const messageSchema = new Schema({
    channelID: {type: String, required: true},
    body: String,
    createdAt: Date,
    creator: {
        email: String,
        userName: String,
        firstName: String,
        lastName: String
    },
    editedAt: Date
});

module.exports = { channelSchema, messageSchema }