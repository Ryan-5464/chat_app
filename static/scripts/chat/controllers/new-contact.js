function addAddContactEventListenerToAddContactInput() {
    const addContactInput = ConfigureAddContactInput();
    addContactInput.addEventListener('keydown', (e) => {
        if (e.key !== 'Enter' || !addContactInput.value.trim()) return;
        console.log("add contact email: ", addContactInput.value)
        addContactInput.__controller.AddContact(addContactInput.value);
        addContactInput.value = ''
    });
};

function ConfigureAddContactInput () {
    const addContactInput = document.getElementById(APP.ID.CONTACT.INPUT.NEW_CONTACT_EMAIL);
    addContactInput.__controller = {
        AddContact: (input) => AddContactHandler(input),
    };
    return addContactInput;
};

