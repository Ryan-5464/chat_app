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
    const addContactInput = document.getElementById('contact-email-input');
    addContactInput.__controller = {
        AddContact: (input) => AddContactHandler(input),
    };
    return addContactInput;
};

