EDIT_USERNAME_ENDPOINT = '/api/profile/name/edit'

document.addEventListener("DOMContentLoaded", function() {
    addEditUserNameListener()

})

function addEditUserNameListener() {
    const userNameInput = document.getElementById('username') 
    userNameInput.addEventListener('keydown', function (e) {
        if (e.key !== 'Enter') return;
        e.preventDefault();
        EditUserNameHandler(userNameInput.value);
    });
};

const EditUserNameRequest = (name) => safeRequest(() => POST(EDIT_USERNAME_ENDPOINT, { Name: name}));

function EditUserNameHandler(name) {
    EditUserNameRequest(name).then(data => {
        console.log("Edit user name response handler: ", data);
        HandleEditUserNameResponse(data);
        const userNameInput = document.getElementById('username')
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
    const currentName = document.getElementById('current-value')
    currentName.innerHTML = name
};
