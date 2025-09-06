
document.addEventListener("DOMContentLoaded", function() {
    addEditUserNameListener();
    GetOnlineStatusHandler();

})

function addEditUserNameListener() {
    const userNameInput = document.getElementById(APP.ID.USER.NAME) 
    userNameInput.addEventListener('keydown', function (e) {
        console.log("edit username")
        if (e.key !== 'Enter') return;
        e.preventDefault();
        EditUserNameHandler(userNameInput.value);
    });
};

const EditUserNameRequest = (name) => safeRequest(() => POST(APP.ENDPOINT.EDIT_USERNAME, { Name: name}));

function EditUserNameHandler(name) {
    EditUserNameRequest(name).then(data => {
        console.log("Edit user name response handler: ", data);
        HandleEditUserNameResponse(data);
        const userNameInput = document.getElementById(APP.ID.USER.NAME)
        userNameInput.value = ''
    }).catch(error => {
        console.error("Edit user name failed => error: ", error);
    });
};

function HandleEditUserNameResponse(data) {
    const callbacks = {
        Name: (data) => RenderUserName(data), 
    };
    return HandleResponse(data, callbacks);
};

function RenderUserName(name) {
    const currentName = document.getElementById(APP.ID.GEN.PFL_USERNAME)
    currentName.innerHTML = name
};
