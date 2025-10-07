<h1>ChatApp - WhatsApp like app.</h1>

<p>Below are short descriptions of the main functionality accompanied by short video demonstrations. Click the images to watch the videos. After that are high level diagrams of the architecture of the app.</p>

<h2>Register User</h2>

[![Watch the video](https://img.youtube.com/vi/MoZRt9O_67c/hqdefault.jpg)](https://www.youtube.com/watch?v=MoZRt9O_67c)

<p>Users can register and login. Password and email verifcation are performed. Users are given an encrypted web token on login/regstration. Users without a token are redirected to the landing page upon trying to access server resources.</p>

<h2>Add Contact</h2>

[![Watch the video](https://img.youtube.com/vi/XLF5jOyDFWg/hqdefault.jpg)](https://www.youtube.com/watch?v=XLF5jOyDFWg)

<p>Users can add contacts by entering the contact's email address.</p>

<h2>Online Status</h2>

[![Watch the video](https://img.youtube.com/vi/X5zlCl32LCw/hqdefault.jpg)](https://www.youtube.com/watch?v=X5zlCl32LCw)

<p>Contacts are able to see each others online status in real time. Users may appear offline by setting their status to stealth mode.</p>

<h2>Create New Chat</h2>

[![Watch the video](https://img.youtube.com/vi/L8zrmb52pho/hqdefault.jpg)](https://www.youtube.com/watch?v=L8zrmb52pho)

<p>Users can create chats by entering a chat name. Chat names are unique to the user only. The user that creates the chat has admin permissions such as editing the name and adding/removing members.</p>

<h2>Add Users to Chat</h2>

[![Watch the video](https://img.youtube.com/vi/0f6sfHE9vwE/hqdefault.jpg)](https://www.youtube.com/watch?v=0f6sfHE9vwE)

<p>Users can be added to the chat by entering their email under the list of chat members.</p>

<h2>Real-Time Messaging</h2>

[![Watch the video](https://img.youtube.com/vi/pisHYtknUJ0/hqdefault.jpg)](https://www.youtube.com/watch?v=pisHYtknUJ0)

<p>Users can see new messages, edits to messages, and message deletions in real-time</p>

<h2>Unread Message Tracking</h2>

[![Watch the video](https://img.youtube.com/vi/piBxBzaGaqo/hqdefault.jpg)](https://www.youtube.com/watch?v=piBxBzaGaqo)

<p>The app will track counts for new unread messages for each user. </p> 

<h2>User Permissions</h2>

[![Watch the video](https://img.youtube.com/vi/cRwuk0wtSlo/hqdefault.jpg)](https://www.youtube.com/watch?v=cRwuk0wtSlo)

<p>Users are verified as the owners of chats and messages using their encrypted web token. Owners have the ability to edit and deletes messages, edit chat names, add/remove users from chats. If the owner of a chat leaves it, a new admin is automatically picked from the remianing members, if no members exist, the chat is deleted.</p>

<h2>Adding Chat Members as New Contacts</h2>

[![Watch the video](https://img.youtube.com/vi/gFTw9vpm8yg/hqdefault.jpg)](https://www.youtube.com/watch?v=gFTw9vpm8yg)

<p>Users are able to add any members of a chat as new contacts.</p>

<h2>ChatApp Architecture</h2>

<img width="1345" height="1000" alt="chatapp-architecture" src="https://github.com/user-attachments/assets/d32087e7-635d-429e-8bc4-0f9120154957" />

<p>The image shows the high level architecture for the app. The design was chosen to minimize coupling of functionality, to clearly seperate data models from domain entities and orchestration from data retrieval, and allow the database layer to be swapped for other non-SQL solutions in the future if needed.</p>

<h2>User Authentication</h2>

<img width="1549" height="638" alt="secretkeygen" src="https://github.com/user-attachments/assets/e9a855b2-0ca1-4f32-b07f-b97023e59d4e" />

<p>The image shows a high level view of user verification using secret keys. Keys are ethemeral and generated periodically to increase server security and old keys stored temporarily to prevent users being logged out unexpectedly.</p>
