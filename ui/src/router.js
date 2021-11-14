
import Page404 from "./pages/page404";
import Login from "./pages/login/Login";
import UserLayout from "./pages/user/UserLayout";
import AdminLayout from "./pages/admin/Layout";

const routes = [
    {
        exact: true,
        path: "/",
        component: UserLayout
    },
    {
        exact: true,
        path: "/p",
        component: UserLayout
    },
    {
        path: "/ad",
        component: AdminLayout
    },
    {
        exact: true,
        path: "/log",
        component: Login
    },
    {
        path: "/",
        component: Page404
    },
];

export default routes;