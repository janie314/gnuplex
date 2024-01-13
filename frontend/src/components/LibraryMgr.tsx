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
      <div className="librarymgr">
        <Tabs>
          <Tab key="dirs" title="Directories" />
          <Tab key="tags" title="Tagging" />
        </Tabs>
        <FilesAndDirectoriesTab selected={true} closeHook={props.closeHook} />
      </div>
    );
  }
  return <></>;
}

export { LibraryMgr };
