import {createStore} from 'redux';


//reducer
function loginManage(state = {}, action) {
    switch (action.type) {
        case "loginIn":
            return {
                username: action.username,
            }
        case "logout":
            return {}
        default:
            return state;
    }
}

let store = createStore(loginManage);

function loginIn(name) {
    store.dispatch({
        type: "loginIn",
        username: name,
    })
}

function logout() {
    store.dispatch({
        type: "logout",
    })
}

export {loginIn, logout};
export default store;