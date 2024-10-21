function WorkingSpinnerTSX(props: { visible: boolean }) {
  return (
    <span className="font-mono text-base m-2">{props.visible ? "↺" : " "}</span>
  );
}

export { WorkingSpinnerTSX };
