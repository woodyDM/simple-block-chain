import {Route, Switch} from 'react-router-dom';
import React from "react";

function NestRoute(props) {
    const routes = props.routes ? props.routes : [];
    return (<Switch>
        {
            routes.map((route, i) => (
                <RouteWrapper key={i} {...route}/>
            ))
        }
    </Switch>);
}

function RouteWrapper(route) {
    return (
        <Route
            path={route.path}
            render={p => <route.component {...p} routes={route.routes}/>}/>
    );
}

export default NestRoute;