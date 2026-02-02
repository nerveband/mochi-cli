import React from "react";

interface TerminalWindowProps {
  children: React.ReactNode;
}

export const TerminalWindow: React.FC<TerminalWindowProps> = ({ children }) => {
  return (
    <div
      style={{
        borderRadius: 24,
        overflow: "hidden",
        boxShadow: "0 20px 50px rgba(0, 0, 0, 0.2)",
        backgroundColor: "#1a1a1a",
        position: "relative",
        border: "1px solid rgba(255, 255, 255, 0.1)",
      }}
    >
      {/* Title bar - macOS style */}
      <div
        style={{
          height: 60,
          backgroundColor: "#2a2a2a",
          borderBottom: "1px solid rgba(255, 255, 255, 0.05)",
          display: "flex",
          alignItems: "center",
          padding: "0 24px",
          position: "relative",
        }}
      >
        {/* Traffic lights */}
        <div style={{ display: "flex", gap: 10 }}>
          <div style={{ width: 14, height: 14, borderRadius: "50%", backgroundColor: "#ff5f57" }} />
          <div style={{ width: 14, height: 14, borderRadius: "50%", backgroundColor: "#febc2e" }} />
          <div style={{ width: 14, height: 14, borderRadius: "50%", backgroundColor: "#28c840" }} />
        </div>

        {/* Title */}
        <div
          style={{
            position: "absolute",
            left: "50%",
            transform: "translateX(-50%)",
            color: "rgba(255, 255, 255, 0.4)",
            fontSize: 16,
            fontWeight: 500,
            letterSpacing: "0.02em",
          }}
        >
          mochi — Terminal
        </div>
      </div>

      {/* Terminal content */}
      <div
        style={{
          padding: "50px 60px",
          fontSize: 32,
          lineHeight: 1.6,
          color: "#f7f7f7",
          minHeight: 700,
          position: "relative",
        }}
      >
        {children}
      </div>
    </div>
  );
};
