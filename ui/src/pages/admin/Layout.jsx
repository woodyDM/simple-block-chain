import React from 'react';
import './layout.css';
import {Layout, Menu} from 'antd';
import {DesktopOutlined, PieChartOutlined} from '@ant-design/icons';
import NestRoute from "../../components/NestRoute";
import routes from './router';
import AdminPageWrapped from "./AdminPageWrapper";

const {Header, Content, Footer, Sider} = Layout;

class MyLayout extends React.Component {
    constructor(props) {
        super(props);
        this.map = {
            '1': '/ad/page1',
            '2': "/ad/page2",
            '3': "/ad/page3",
        }
    }

    state = {
        collapsed: false,
    };

    handleClick = e => {
        const key = e.key;
        if (this.map[key]) {
            this.props.history.push(this.map[key]);
        }
    }

    onCollapse = collapsed => {
        this.setState({collapsed});
    };

    render() {
        const {collapsed} = this.state;
        return (
            <AdminPageWrapped>
                <div id="components-layout-demo-side">
                    <Layout style={{minHeight: '100vh'}}>
                        <Sider collapsible collapsed={collapsed} onCollapse={this.onCollapse}>
                            <div className="logo"/>
                            <Menu
                                onClick={this.handleClick}
                                theme="dark"
                                defaultSelectedKeys={['1']}
                                mode="inline">
                                <Menu.Item key="1" icon={<DesktopOutlined/>}>
                                    图片上传
                                </Menu.Item>
                                <Menu.Item key="2" icon={<DesktopOutlined/>}>
                                    图片管理
                                </Menu.Item>
                                <Menu.Item key="3" icon={<PieChartOutlined/>}>
                                    数据管理
                                </Menu.Item>
                            </Menu>
                        </Sider>
                        <Layout className="site-layout">
                            <Header className="site-layout-background" style={{padding: 0}}/>
                            <Content style={{margin: '0 16px'}}>
                                <div className="site-layout-background"
                                     style={{margin: '16px 0', padding: 24, minHeight: 360}}>
                                    <NestRoute routes={routes}/>
                                </div>
                            </Content>
                            <Footer style={{textAlign: 'center'}}>shadiaotu ©2021 Created by wd</Footer>
                        </Layout>
                    </Layout>
                </div>
            </AdminPageWrapped>
        );
    }
}

export default MyLayout;