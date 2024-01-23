import "../App.css";
import "./LibraryMgr.css";
import { FilesAndDirectoriesTab } from "./LibraryMgr/FilesAndDirectoriesTab.tsx";
import { TaggingTab } from "./LibraryMgr/TaggingTab.tsx";
import { Tab, Tabs } from "@nextui-org/react";
import { useState } from "react";

function LibraryMgr(props: {
  visible: boolean;
  medialist: string[];
  setMediadirInputPopup: React.Dispatch<React.SetStateAction<boolean>>;
  closeHook: () => void;
}) {
  const [selected, setSelected] = useState("dirs");
  if (props.visible) {
    return (
      <div className="librarymgr flex flex-col shadow z-1 w-3/4 p-4 space-y-3 bg-slate-100">
        <Tabs
          color="primary"
          variant="bordered"
        >
          <Tab
            key="dirs"
            title="Directories"
          >
            <FilesAndDirectoriesTab
              closeHook={props.closeHook}
            />
          </Tab>
          <Tab key="tags" title="Tagging">
            <TaggingTab
              closeHook={props.closeHook}
              medialist={props.medialist}
            >
            </TaggingTab>
          </Tab>
        </Tabs>
      </div>
    );
  }
  return <></>;
}

export { LibraryMgr };
