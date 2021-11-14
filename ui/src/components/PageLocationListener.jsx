import {withRouter} from "react-router";
import {useEffect} from "react";

function PageLocationListener(props) {
    const cur = props.location.pathname + props.location.search;
    useEffect(() => {
        console.log("Do something with current path :" + cur);
    }, [cur]);
    return (
        <div></div>
    );
}

export default withRouter(PageLocationListener);