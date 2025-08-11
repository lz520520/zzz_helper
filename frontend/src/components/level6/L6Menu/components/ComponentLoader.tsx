import React, { ReactElement, useEffect, useState } from 'react';

type ComponentLoaderProps = {
    componentFn: React.FC<L6MenuComponentProps>;
    menuProps: L6MenuComponentProps;
};
export type L6MenuComponentProps = {
    text: string;
    open: boolean;
    close: () => void;
    rowArgs: Record<string, any> | string;
    callbackEvent: (extParams: Record<string, any>) => Promise<void>; // 回调
    submit?: () => void,
};

const ComponentLoader: React.FC<ComponentLoaderProps> = ({ componentFn, menuProps }) => {
    const [component, setComponent] = useState<ReactElement<any, any> | null>(null);

    useEffect(() => {
        const loadComponent = async () => {
            try {
                // 直接调用传入的函数组件定义，获取组件
                const DynamicComponent = componentFn;
                // 设置组件状态
                setComponent(<DynamicComponent {...menuProps} />);
            } catch (error) {
                console.error('Error loading component:', error);
            }
        };

        // 调用加载组件函数
        loadComponent();
    }, [componentFn, menuProps]);

    return component;
};

export default ComponentLoader;
