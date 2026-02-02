import React from "react";

interface OutputDisplayProps {
  type: "json" | "table" | "config" | "success";
  fadeInProgress: number;
}

const COLORS = {
  key: "#badd77",
  string: "#a7ce5f",
  number: "#febc2e",
  bracket: "#6bcbff",
  success: "#28c840",
};

export const OutputDisplay: React.FC<OutputDisplayProps> = ({
  type,
  fadeInProgress,
}) => {
  return (
    <div
      style={{
        opacity: fadeInProgress,
        transform: `translateY(${(1 - fadeInProgress) * 10}px)`,
      }}
    >
      {type === "json" && <JsonOutput />}
      {type === "table" && <TableOutput />}
      {type === "config" && <ConfigOutput />}
      {type === "success" && <SuccessOutput />}
    </div>
  );
};

const JsonOutput: React.FC = () => (
  <pre
    style={{
      margin: 0,
      fontFamily: "inherit",
      fontSize: 28,
      lineHeight: 1.6,
    }}
  >
    <span style={{ color: COLORS.bracket }}>{"["}</span>
    {"\n  "}
    <span style={{ color: COLORS.bracket }}>{"{"}</span>
    {"\n    "}
    <span style={{ color: COLORS.key }}>"id"</span>
    <span style={{ color: "#f7f7f7" }}>: </span>
    <span style={{ color: COLORS.string }}>"spanish-vocab"</span>
    <span style={{ color: "#f7f7f7" }}>,</span>
    {"\n    "}
    <span style={{ color: COLORS.key }}>"name"</span>
    <span style={{ color: "#f7f7f7" }}>: </span>
    <span style={{ color: COLORS.string }}>"Spanish Vocabulary"</span>
    {"\n  "}
    <span style={{ color: COLORS.bracket }}>{"}"}</span>
    {"\n"}
    <span style={{ color: COLORS.bracket }}>{"]"}</span>
  </pre>
);

const TableOutput: React.FC = () => (
  <div style={{ fontFamily: "inherit", fontSize: 28 }}>
    <div style={{ color: COLORS.key, marginBottom: 16, fontWeight: 600 }}>
      {"ID          CONTENT           UPDATED"}
    </div>
    <div style={{ color: "#f7f7f7", borderTop: "1px solid rgba(255,255,255,0.1)", paddingTop: 16 }}>
      <div style={{ marginBottom: 12 }}>
        <span style={{ color: COLORS.string }}>a1b2c3d  </span>
        <span>Hola / Hello      </span>
        <span style={{ color: "rgba(255,255,255,0.6)" }}>2m ago</span>
      </div>
      <div style={{ marginBottom: 12 }}>
        <span style={{ color: COLORS.string }}>e4f5g6h  </span>
        <span>Gracias / Thanks  </span>
        <span style={{ color: "rgba(255,255,255,0.6)" }}>5m ago</span>
      </div>
      <div style={{ marginBottom: 12 }}>
        <span style={{ color: COLORS.string }}>i7j8k9l  </span>
        <span>Buenos días / Goo </span>
        <span style={{ color: "rgba(255,255,255,0.6)" }}>10m ago</span>
      </div>
    </div>
  </div>
);

const ConfigOutput: React.FC = () => (
  <div style={{ fontFamily: "inherit", fontSize: 28 }}>
    <div style={{ marginBottom: 24 }}>
      <span style={{ color: COLORS.success, marginRight: 16 }}>*</span>
      <span style={{ color: COLORS.key, width: 160, display: "inline-block" }}>default</span>
      <span style={{ color: "rgba(255,255,255,0.7)" }}>35c7bc879de3fb146ad0b0e9</span>
    </div>
    <div style={{ marginTop: 24, color: "rgba(255,255,255,0.5)", fontSize: 24 }}>
      <span style={{ color: COLORS.success }}>*</span> = active profile
    </div>
  </div>
);

const SuccessOutput: React.FC = () => (
  <div style={{ color: COLORS.success, fontSize: 32, fontWeight: 500 }}>
    ✓ Card created successfully
  </div>
);
