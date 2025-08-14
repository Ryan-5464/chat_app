const BASEURL = 'http://localhost:8081'
const DEL_MSG_ENDPOINT = '/api/message/delete'
const EDIT_CHAT_NAME_ENDPOINT = '/api/chat/edit'

function POST(json) {
    return {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(json),
    }
}

function DELETE() {
    return {
        method: 'DELETE',
    }
}