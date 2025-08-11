import { L6MenuComponentProps } from '@/components/level6/L6Menu';
import { ExclamationCircleFilled } from '@ant-design/icons';
import { Modal } from 'antd';
import React from 'react';

const { confirm } = Modal;
const showConfirm = () => {
    confirm({
        title: 'Do you Want to delete these items?',
        icon: <ExclamationCircleFilled />,
        content: 'Some descriptions',
        onOk() {
            console.log('OK');
        },
        onCancel() {
            console.log('Cancel');
        },
    });
};

const TwiceConfirmDialog: React.FC<L6MenuComponentProps> = (props) => {
    return (
        <Modal
            title={
                <>
                    <ExclamationCircleFilled />
                    你确认要进行"${props.text}"操作吗?
                </>
            }
            open={props.open}
            onCancel={props.close}
            onOk={async function () {
                await props.callbackEvent({});
                props.close();
                showConfirm();
            }}
        >
            <ExclamationCircleFilled />
        </Modal>
    );
};

export default TwiceConfirmDialog;
