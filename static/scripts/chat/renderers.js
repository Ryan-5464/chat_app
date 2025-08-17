
class Renderer {
    configs = {
        Chats: {
            containerId: 'chats-container',
            elemFactory: ChatElement,
        },
        Messages: {
            containerId: 'messages-container',
            elemFactory: MessageElement,
        },
        Contacts: {
            containerId: 'contacts-container',
            elemFactory: ContactElement,
        }
    }

    addConfig(key, config) {
        configs[key] = config
    }
    
    render(configSelector, data, overwrite) {
        const config = configs[configSelector]
        const container = document.getElementById(config.containerId)
        if (overwrite == true) {
            container.innerHTML = ''
        }
        
        Object.values(data).forEach(obj => {
            container.appendChild(elemFactory(obj))
        })
    }

}
