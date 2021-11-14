import {connect} from "react-redux";
import {useEffect} from "react";
import {loginIn} from "../../store/loginState";

function AdminPageWrapper(props) {
    const username = props.loginUser;
    const {children} = props;
    const getUser = () => Promise.resolve("adminName");

    useEffect(() => {
        if (!username) {
            getUser().then(name => {
                console.log("loginIn:" + name);
                loginIn(name);
            });
        }
    }, [username]);
    return (
        <div>
            {username ? children : ""}
        </div>
    );
}

function mapStateToProps(state) {
    return {
        loginUser: state.username,
    }
}

const mapDispatchToProps = dispatch => {
    return {}
}
const AdminPageWrapped = connect(mapStateToProps, mapDispatchToProps)(AdminPageWrapper);

export default AdminPageWrapped;