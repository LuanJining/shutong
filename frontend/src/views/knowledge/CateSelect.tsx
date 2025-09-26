import "./styles/cate-select.scss"
import { Modal } from "antd";
import { RightOutlined } from "@ant-design/icons";

interface IProps {
    open: boolean,
    setOpen: Function
}

export default function CateSelect({ open, setOpen }: IProps) {
    return (
        <Modal
            open={open}
            centered
            title="选择类别"
            className="cate-modal"
            width='50%'
            onCancel={() => { setOpen(false) }}
        >
            <div className="cate-content flex mgT24">
                <div className="cate-box flex space-between">
                    <div className="cate-inner">
                        <div className="cate-item">部门知识库</div>
                        <div className="cate-item">企业知识库</div>
                    </div>
                    <div className="flex-center"><RightOutlined className="lg-fs pointer" /></div>
                </div>
                <div className="cate-box flex space-between">
                    <div className="cate-inner">
                        <div className="cate-item">部门知识库</div>
                        <div className="cate-item">企业知识库</div>
                    </div>
                    <div className="flex-center"><RightOutlined className="lg-fs pointer" /></div>
                </div>
                <div className="cate-box">
                    <div className="cate-inner">
                        <div className="cate-item">部门知识库</div>
                        <div className="cate-item">企业知识库</div>
                    </div>
                </div>
            </div>
        </Modal>
    )
}
