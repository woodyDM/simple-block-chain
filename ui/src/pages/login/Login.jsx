import {Button, Form, Input, Layout} from 'antd';
import './Login.css';
import {useHistory} from 'react-router-dom';
import apis from "../../utils/apis";

const {Header, Content} = Layout;

const layout = {
    labelCol: {span: 8},
    wrapperCol: {span: 8},
};
const tailLayout = {
    wrapperCol: {offset: 8, span: 8},
};

const Login = function (props) {

    const history = useHistory();

    const onFinish = values => {
        console.log(values);
        history.push(apis.page.adminHome)
    };

    return (
        <Layout className="site-layout">
            <Header className="site-layout-background" style={{padding: 0}}/>
            <Content style={{margin: '200px 0 200px'}}>
                <Form
                    {...layout}
                    name="basic"
                    onFinish={onFinish}
                >
                    <Form.Item
                        label="Username"
                        name="Name"
                        rules={[{required: true, message: 'Please input your username!!!'}]}
                    >
                        <Input/>
                    </Form.Item>
                    <Form.Item
                        label="Password"
                        name="Pass"
                        rules={[{required: true, message: 'Please input your password!!!'}]}
                    >
                        <Input.Password/>
                    </Form.Item>
                    <Form.Item {...tailLayout}>
                        <Button type="primary" htmlType="submit">
                            Login
                        </Button>
                    </Form.Item>
                </Form>
            </Content>
        </Layout>
    );
}

export default Login;
