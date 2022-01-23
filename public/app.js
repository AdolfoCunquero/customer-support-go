var app = new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        email: null, // Email address used for grabbing an avatar
        username: "61e12342d9388a555a404379", // Our username
        joined: false, // True if email and username have been filled in
        contacts:[],
        messages:[]
    },

    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            console.log(msg);

            if(msg.type == "joined"){
                self.contacts = msg.contactsUUID;
            }
            else if(msg.type ==  "connectContact"){
                self.contacts.push({uuid:msg.uuid})
            }
            else if (msg.type == "disconectContact"){
                self.contacts = self.contacts.filter(c => c.uuid != msg.uuid)
            }
            else if (msg.type==""){
                // self.chatContent += '<div class="chip">'
                //     + '<img src="' + self.gravatarURL(msg.email) + '">' // Avatar
                //     + msg.username
                // + '</div>'
                // + emojione.toImage(msg.message) + '<br/>'; // Parse emojis

                self.messages.push(msg)

                var element = document.getElementById('chat-messages');
                element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
            }
        });
    },

    methods: {
        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        email: this.email,
                        username: this.username,
                        clientId:"61e12342d9388a555a404379",
                        messageContent:{
                            messageType:1,
                            text: $('<p>').html(this.newMsg).text() // Strip out html
                        }
                        
                    }
                ));
                this.newMsg = ''; // Reset newMsg
            }
        },

        join: function () {
            if (!this.email) {
                Materialize.toast('You must enter an email', 2000);
                return
            }
            if (!this.username) {
                Materialize.toast('You must choose a username', 2000);
                return
            }
            this.email = $('<p>').html(this.email).text();
            this.username = $('<p>').html(this.username).text();
            this.joined = true;
        },

        gravatarURL: function(email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
        }
    }
});