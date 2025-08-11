import React from "react";
import {Row} from "antd";

const L6Row: React.FC<{children: React.ReactNode,marginBottom?:string|number,  width?: string|number, height?: string| number}> = ({children,marginBottom, width, height}) => {
    return (
        <Row gutter={16}  style={{marginBottom: marginBottom?marginBottom:8, width: width, height: height}}>
            {children}
        </Row>
    );
};

export default  L6Row
