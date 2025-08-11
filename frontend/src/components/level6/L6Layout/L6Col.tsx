import React, {CSSProperties} from "react";
import {Col} from "antd";

const L6Col: React.FC<{span: number|string,paddingRight?:number|string, style?:CSSProperties, children: React.ReactNode, flex?: number}> = ({span,paddingRight,style, children, flex}) => {
    return (
        <Col span={span} style={{paddingRight: paddingRight?paddingRight:8, ...style}} flex={flex}>
            {children}
        </Col>
    );
};

export default  L6Col
