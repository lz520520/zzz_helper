import PageWrapTabs from "@/pages/components/PageWrapTabs";
import DBTable from "@/pages/components/DBTable/DBTable";


const Content: React.FC<{id?: string}> = ({id}) => {
    return (<DBTable moduleName={"driver_fuzz"} keyName={"id"} operation={{}}/>)
}

const DriverFuzzHistoryPage: React.FC =() => {
    return (
        <PageWrapTabs
            tag={"driver_parser"}
            content={(id?: string)=> <Content id={id}/>}
        />
    );
}


export default DriverFuzzHistoryPage;
