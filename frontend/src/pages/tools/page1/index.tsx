import PageWrapTabs from "@/pages/components/PageWrapTabs";
import DBTable from "@/pages/components/DBTable/DBTable";

const Content: React.FC<{id?: string}> = ({id}) => {
    return (<DBTable moduleName={"ip_pool"} keyName={"id"} operation={{
        Create: true,
    }}/>)
}

const Page1: React.FC =() => {
    return (
        <PageWrapTabs
            tag={"page1"}
            content={(id?: string)=> <Content id={id}/>}
        />
    );
}


export default Page1;
