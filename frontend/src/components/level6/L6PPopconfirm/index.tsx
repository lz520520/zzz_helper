import React, { useEffect, useRef, useState} from "react";
import {message, Popconfirm} from "antd";


export type ClientPosition  = {
    clientX: number,
    clientY: number
}

const usePopconfirm = () => {
    const [open, setOpen] = useState(false);
    const [currentEvent, setCurrentEvent] = useState<(() => void)|null>(null);
    const [title, setTitle] = useState('');
    const [position, setPosition] = useState<ClientPosition>({clientX: 0, clientY: 0});
    const popconfirmRef = useRef<HTMLDivElement>(null); // 创建一个 ref 来引用 Popconfirm

    const showPopconfirm = (newTitle: string, eventFn: () => void, pos: ClientPosition) => {
        setCurrentEvent(() => eventFn); // 保存要执行的事件
        setTitle(newTitle);
        setOpen(true);
        setPosition(pos); // 设置鼠标点击的位置

    };

    const handleOk = () => {
        if (currentEvent) {
            currentEvent(); // 执行保存的事件并传递参数
        }
        setOpen(false); // 关闭 Popconfirm
    };
    const handleCancel = () => {
        setOpen(false);
    };
    // // 监听点击事件以关闭 Popconfirm
    // useEffect(() => {
    //     const handleClickOutside = (event: MouseEvent) => {
    //         if (popconfirmRef.current && !popconfirmRef.current.contains(event.target as Node)) {
    //             setOpen(false);
    //         }
    //     };
    //     if (open) {
    //         document.addEventListener('mouseup', handleClickOutside, { once: true });
    //     }
    //
    //     return () => {
    //         document.removeEventListener('mouseup', handleClickOutside);
    //     };
    // }, [open]);

    const popconfirm = (
        <div
            ref={popconfirmRef}
            style={{
                position: 'fixed',
                top: `${position.clientY}px`,
                left: `${position.clientX}px`,
                zIndex: 9999,
            }}
        >
            <Popconfirm
                title={`你确认要进行"${title}"操作吗?`}
                open={open}
                onConfirm={handleOk}
                onCancel={handleCancel}
                okText="确认"
                cancelText="取消"
            >
                <div  style={{ display: 'none' }} /> {/* 隐藏触发按钮 */}
            </Popconfirm>
        </div>
    );

    return { showPopconfirm, popconfirm };
};

export default usePopconfirm;
