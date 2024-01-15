import "../App.css";
import "./LibraryMgr.css";
import { FilesAndDirectoriesTab } from "./LibraryMgr/FilesAndDirectoriesTab.tsx";
import { Tab, Tabs } from "@nextui-org/react";

function LibraryMgr(props: {
  visible: boolean;
  setMediadirInputPopup: React.Dispatch<React.SetStateAction<boolean>>;
  closeHook: () => void;
}) {
  if (props.visible) {
    return (
      <div className="librarymgr flex flex-col shadow z-1 w-3/4 p-4 space-y-3 bg-med-blue">
        <Tabs color="primary">
          <Tab key="dirs" title="Directories" />
          <Tab key="tags" title="Tagging" />
        </Tabs>
        <FilesAndDirectoriesTab
          selected={true}
          closeHook={props.closeHook}
        />
      </div>
    );
  }
  return <></>;
}

export { LibraryMgr };
