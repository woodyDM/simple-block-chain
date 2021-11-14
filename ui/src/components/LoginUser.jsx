import {connect} from "react-redux";
import {logout} from "../store/loginState";
import {WrapAdmin} from "../utils/request";
import axios from "axios";
import apis from "../utils/apis";
import {useHistory} from 'react-router-dom';

function LoginUser(props) {
    const his = useHistory();
    const btClick = () => {
        WrapAdmin(his,() => axios.post(apis.api.logout))().then(d =>
            logout()
        )
    }

    return (
        <div>LoginUser:{props.loginUser}
            <button onClick={btClick}>Logout</button>
        </div>
    );
}

function mapPros(state) {
    return {
        loginUser: state.username,
    }
}

function mapDispatch(dispatch) {
    return {}
}

export default connect(mapPros, mapDispatch)(LoginUser);