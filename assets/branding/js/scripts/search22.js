$(document).ready(function () {
    // Prevent form submission on Enter keypress
    $("#userSearchForm, #userSearchForm2").on("submit", function (event) {
        event.preventDefault(); // Prevent the form from submitting and reloading the page
    });

    // Handle input event for real-time search
    $("#searchInput").on("input", function () {
        performSearch(); // Call the search function
    });

    // Handle the click event on the search button
    $("#searchButton").on("click", function () {
        performSearch(); // Call the search function
    });

    // Perform search for user contacts
    function performSearch() {
        var username = $("#searchInput").val(); // Get the value from the input field

        $.ajax({
            type: "POST",
            url: "/api/functions/search",
            data: { username: username },
            success: function (users) {
                var userList = $("#chatContactTab"); // Assuming you have a <ul> or <div> with id="chatContactTab"
                userList.empty(); // Clear previous results

                if (users.length === 0) {
                    userList.append("<li>No users found.</li>");
                } else {
                    $.each(users, function (index, user) {
                        var userItem = `
                            <li class="contacts-item friends">
                                <a class="contacts-link" href="/dash/chats?userID=${user.ID}">
                                    <div class="avatar avatar-online">
                                        <img src="./..//_assets/media/avatar/${user.Avatar}" alt="">
                                    </div>
                                    <div class="contacts-content">
                                        <div class="contacts-info">
                                            <h6 class="chat-name text-truncate">${user.Username}</h6>
                                            <div class="chat-time">${user.Status}</div>
                                        </div>
                                        <div class="contacts-texts">
                                            <p class="text-truncate">n/a</p>
                                        </div>
                                    </div>
                                </a>
                            </li>
                        `;
                        userList.append(userItem);

                        if (user.chatID) {
                            setCookie('chatID', user.chatID, 1); // Cookie expires in 1 day
                        }
                    });
                }
            },
            error: function (xhr, status, error) {
                console.error('Error performing user search:', error);
            }
        });
    }

    // Handle input event for real-time search in StartChat
    $("#searchInput2").on("input", function () {
        performSearch2(); // Call the search function
    });

    // Handle the click event on the search button for StartChat
    $("#searchButton2").on("click", function () {
        performSearch2(); // Call the search function
    });

    // Perform search for StartChat
    function performSearch2() {
        var username = $("#searchInput2").val(); // Corrected to use the right input field ID

        $.ajax({
            type: "POST",
            url: "/api/functions/search",
            data: { username: username },
            success: function (users) {
                var userList = $("#contactsList"); // Assuming you have a <ul> or <div> with id="contactsList"
                userList.empty(); // Clear previous results

                if (users.length === 0) {
                    userList.append("<li>No users found.</li>");
                } else {
                    $.each(users, function (index, user) {
                        var userItem = `
                            <li class="list-group-item">
                                <a class="media" href="/dash/chats?userID=${user.ID}">
                                    <div class="avatar avatar-offline mr-2">
                                        <img src="./..//_assets/media/avatar/${user.Avatar}" alt="">
                                    </div>
                                    <div class="media-body">
                                        <h6 class="text-truncate">
                                            <a href="#" class="text-reset">${user.Username}</a>
                                        </h6>
                                            <p class="text-muted mb-0">${user.Status}</p>
                                    </div>
                                </a>
                            </li>
                        `;
                        userList.append(userItem);


                        if (user.chatID) {
                            setCookie('chatID', user.chatID, 1); // Cookie expires in 1 day
                        }
                    });
                }
            },
            error: function (xhr, status, error) {
                console.error('Error performing StartChat search:', error);
            }
        });
    }

    // Load recent chats on page load
    RecentChats();

    function RecentChats() {
        $.ajax({
            type: "GET",
            url: "/api/functions/recent",
            success: function(chats) {
                var chatList = $("#chatContactTab"); // Replace with your container ID
                chatList.empty(); // Clear previous results

                if (chats.length === 0) {
                    chatList.append("<li>No recent chats found.</li>");
                } else {
                    $.each(chats, function(index, chat) {
                        var chatItem = `
                            <li class="contacts-item friends">
                                <a class="contacts-link" href="/dash/chats?userID=${chat.peopleID}">
                                    <div class="avatar avatar-online">
                                        <img src="/_assets/media/2.png" alt="">
                                    </div>
                                    <div class="contacts-content">
                                        <div class="contacts-info">
                                            <h6 class="chat-name text-truncate">${chat.name}</h6>
                                            <div class="chat-time">${chat.createdAt}</div>
                                        </div>
                                        <div class="contacts-texts">
                                            <p class="text-truncate">n/a</p>
                                        </div>
                                    </div>
                                </a>
                            </li>
                        `;
                        chatList.append(chatItem);

                        if (chat.chatID) {
                            setCookie('chatID', chat.chatID, 1); // Cookie expires in 1 day
                        }
                    });
                }
            },
            error: function(xhr, status, error) {
                console.error('Error fetching recent chats:', error);
            }
        });
    }

    // Function to set a cookie
    function setCookie(name, value, days) {
        var expires = "";
        if (days) {
            var date = new Date();
            date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
            expires = "; expires=" + date.toUTCString();
        }
        document.cookie = name + "=" + (value || "") + expires + "; path=/";
    }
});
