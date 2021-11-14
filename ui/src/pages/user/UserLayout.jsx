import PageLocationListener from "../../components/PageLocationListener";
import './layout.css';
import {withRouter} from "react-router";
import Header from "../common/Header";
import Footer from "../common/Footer";
import React from "react";
import Home from "./Home";

function UserLayout(props) {
    return (
        <div>
            <PageLocationListener/>
            <Header/>
            <div id="Wrapper">
                <div id="Content">
                    <Home/>
                </div>
            </div>
            <Footer/>
        </div>
    );
}

export default withRouter(UserLayout);