import Page1 from './Page1';
import Page2 from './Page2';
import Page3 from './Page3';
import Page404 from '../page404';

const routes = [
    {
        exact: true,
        path: "/ad",
        component: Page1
    },
    {
        exact: true,
        path: "/ad/page1",
        component: Page1
    },
    {
        exact: true,
        path: "/ad/page2",
        component: Page2
    },
    {
        exact: true,
        path: "/ad/page3",
        component: Page3
    },
    {
        path: "/ad",
        component: Page404
    },
];

export default routes;