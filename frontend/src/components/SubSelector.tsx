import { useState } from "react";
import { API, type SubTrack } from "../lib/API";

function SubSelector(props: { subs: SubTrack[] | null }) {
  const selectedID =
    (props.subs || []).filter((sub) => sub.selected)[0]?.id || -1;
  const [subID, setSubID] = useState(selectedID);
  return (
    <select
      className="select select-sm select-bordered ml-1 mr-1"
      value={subID}
      onChange={(e) => {
        const newID = Number(e.target.value);
        setSubID(newID);
        API.setSubTrack(newID);
      }}
    >
      <option value={-1}>Off</option>
      {(props.subs || []).map((sub) => (
        <option key={sub.id} value={sub.id}>
          {sub.id}
        </option>
      ))}
    </select>
  );
}

export { SubSelector };
