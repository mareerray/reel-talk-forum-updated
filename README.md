*REAL-TIME-FORUM*

<h2>Objectives</h2>
<h4>On this project you will have tofocus on a few points:</h4>
<li>Registration and Login</li>
<li>Creation of posts</li>
    <ul>
    <li>Commenting posts</li>
    </ul>
<li>Private Messages</li>
<br>
<h2>Registration and Login</h2>

<h4>To be able to use the new and upgraded forum users will have to register and login. This is premium stuff. The registration and login process should take in consideration the following features:</h4>
<br>
<li>Users must be able to fill a register form to register into the forum. They will have to provide at least:</li>
    <ul>
        <li>Nickname</li>
        <li>Age</li>
        <li>Gender</li>
        <li>First Name</li>
        <li>Last Name</li>
        <li>E-mail</li>
        <li>Password</li>
    </ul>
<li>The user must be able to connect using either the nickname or the e-mail combined with the password.</li>
<li>The user must be able to log out from any page on the forum.</li>
<br>
<h2>Posts and Comments</h2>
<h4>This part is pretty similar to the first forum. Users must be able to:</h4>
<li>Create posts</li>
    <ul>
        <li>Posts will have categories</li>
    </ul>
<li>Crate comments on the posts</li>
<li>See posts in the feed display</li>
    <ul>
        <li>See comments only if they click on a post</li>    
    </ul>
<br>
<h2>Private Massages</h2>
<h4>Users will be able to send private messages to each other, so you will need to create a chat, where it will exist :</h4>
<li>A section to show who is online/offline and able to talk to:</li>
    <ul>
        <li>This section must be organized by the last message sent (just like discord). If the user is new and does not present messages you must organize it in alphabetic order.
        </li>
        <li>The user must be able to send private messages to the users who are online. 
        </li>
        <li>This section must be visible at all times.
        </li>
    </ul>
<li>A section that when clicked on the user that you want to send a message, reloads the past messages. Chats between users must:</li>
    <ul>
    <li>Be visible, for this you will have to be able to see the previous messages that you had with the user</li>
    <li>Reload the last 10 messages and when scrolled up to see more messages you must provide the user with 10 more, without spamming the scroll event. Do not forget what you learned!! (Throttle, Debounce)</li>
    </ul>
<li>Messages must have a specific format:</li>
    <ul>
    <li>A date that shows when the message was sent</li>
    <li>The user name, that identifies the user that sent the message</li>
    </ul>
<p>As it is expected, the messages should work in real time, in other words, if a user sends a message, the other user should receive the notification of the new message without refreshing the page. Again this is possible through the usage of WebSockets in backend and frontend.</p>

<h1>Allowed Packages</h1>
<br>
<li>All standard go packages are allowed</li>
<li>Gorilla websocket</li>
<li>sqlite3</li>
<li>bcrypt</li>
<li>UUID</li>