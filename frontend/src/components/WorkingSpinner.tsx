function WorkingSpinnerTSX(props: { visible: boolean }) {
  if (props.visible) {
    return <span style={{ font: "16px monospace;", margin: "2px" }}>↺</span>;
  }
  return <span style={{ font: "16px monospace;", margin: "2px" }}></span>;
}

export { WorkingSpinnerTSX };
