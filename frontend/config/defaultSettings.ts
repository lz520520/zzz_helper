import { ProLayoutProps } from '@ant-design/pro-components';


let title = ""
let headerRender: boolean|undefined = undefined

title = "tools"


/**
 * @name
 */
const Settings: ProLayoutProps & {
    pwa?: boolean;
    logo?: string;
} = {
    navTheme: 'light',
    // 拂晓蓝
    colorPrimary: '#1890ff',
    layout: 'side',
    splitMenus: true,

    contentWidth: 'Fluid',
    fixedHeader: false,
    fixSiderbar: true,
    colorWeak: false,
    title: title,
    pwa: false,
    logo: '/favicon.ico',
    iconfontUrl: '',
    siderWidth: 175,
    menu: {
        locale: true, // 设置后，菜单以及tab标签页都可国际化
        defaultOpenAll: true,
        ignoreFlatMenu: true,
    },
    openKeys: false, // 设置后菜单不会自动折叠
    token: {
        header: {
            heightLayoutHeader: 0,
        },
        sider: {

            colorTextMenuSelected: '#0958d9',
            colorBgMenuItemSelected: '#e6f4ff',

            colorTextMenuItemHover: '#0958d9',
            colorBgMenuItemHover: '#e6f4ff',

        },
        // 参见ts声明，demo 见文档，通过token 修改样式
        //https://procomponents.ant.design/components/layout#%E9%80%9A%E8%BF%87-token-%E4%BF%AE%E6%94%B9%E6%A0%B7%E5%BC%8F
    },
    menuHeaderRender: false,
    headerRender: headerRender,

};

export default Settings;
