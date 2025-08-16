
const rendererRegistry = {
    chats: {
        containerId: 'chats-container',
        elemFactory: ChatElement,
    },
    messages: {
        containerId: 'messages-container',
        elemFactory: MessageElement,
    },
    contacts: {
        containerId: 'contacts-container',
        elemFactory: ContactElement,
    }


}

class Renderer {
    constructor() {
        this.config = rendererRegistry 
    }
    
    render(configSelector, data, overwrite) {
        const config = this.config[configSelector]
        const container = document.getElementById(config.containerId)
        if (overwrite == true) {
            container.innerHTML = ''
        }
        
        Object.values(data).forEach(obj => {
            container.appendChild(elemFactory(obj))
        })
    }

}
