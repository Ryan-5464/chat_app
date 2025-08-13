function DeleteMessageRequest() {
    
    const delBtn = document.getElementById('msg-del-btn')
    const chatId = delBtn.getAttribute('data-chatid')
    const messageId = delBtn.getAttribute('data-messageid')
    const userId = delBtn.getAttribute('data-userid')
    
    const params = {
            "chatid": chatId,
            "messageid": messageId,
            "userid": userId,
        }
    const url = BuildURLWithParams(DEL_MSG_ENDPOINT, params)
    
    fetch(url, { method: "DELETE" })
    .then(response => {
        if (!response.ok) throw new Error("Network response was not ok");
        return response.json();
    })
    .then(data => {
        console.log("Delete message response data: ", data)
        renderMessages(data.Messages, true)
    })
    .catch(error => {
        console.log('Error:', error);
    });
}

