import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import 'antd/dist/antd.css'
import {BrowserRouter} from "react-router-dom";
import NestRoute from "./components/NestRoute";
import routes from "./router";
import store from "./store/loginState";
import {Provider} from "react-redux";

ReactDOM.render(
    <React.StrictMode>
        <Provider store={store}>
            <BrowserRouter>
                <div id="main">
                    <NestRoute routes={routes}/>
                </div>
            </BrowserRouter>
        </Provider>
    </React.StrictMode>
    ,
    document.getElementById('root')
);
