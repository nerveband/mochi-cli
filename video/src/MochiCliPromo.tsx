import React from "react";
import {
  AbsoluteFill,
  useCurrentFrame,
  useVideoConfig,
  interpolate,
} from "remotion";
import { TerminalWindow } from "./components/TerminalWindow";
import { TypewriterText } from "./components/TypewriterText";
import { OutputDisplay } from "./components/OutputDisplay";

// Color palette
const COLORS = {
  background: "#f7f7f7",
  text: "#333333",
  accent1: "#a7ce5f",
  accent2: "#badd77",
  button: "#37352f",
  command: "#a7ce5f",
};

// Demo sequences
const SEQUENCES = [
  {
    command: "mochi deck list",
    output: "json",
    description: "List all decks",
  },
  {
    command: "mochi card list --deck spanish-vocab --format table",
    output: "table",
    description: "View cards in any format",
  },
  {
    command: 'mochi card create --deck spanish-vocab --content "# Hola\\n\\nHello"',
    output: "success",
    description: "Create cards instantly",
  },
  {
    command: "mochi deck export vocab.mochi",
    output: "config",
    description: "Import & export decks",
  },
];

export const MochiCliPromo: React.FC = () => {
  const frame = useCurrentFrame();
  const { fps, durationInFrames } = useVideoConfig();

  // Each sequence timing
  const framesPerSequence = Math.floor(durationInFrames / SEQUENCES.length);
  const currentSequenceIndex = Math.floor(frame / framesPerSequence) % SEQUENCES.length;
  const frameInSequence = frame % framesPerSequence;

  const currentSequence = SEQUENCES[currentSequenceIndex];

  // Animation timing
  const typewriterDuration = 30;
  const outputDelay = 40;
  const outputFadeIn = 12;

  return (
    <AbsoluteFill
      style={{
        backgroundColor: COLORS.background,
        fontFamily: "Inter, -apple-system, system-ui, sans-serif",
      }}
    >
      {/* Title */}
      <div
        style={{
          position: "absolute",
          top: 60,
          left: 0,
          right: 0,
          textAlign: "center",
          zIndex: 10,
        }}
      >
        <h1
          style={{
            fontSize: 64,
            fontWeight: 800,
            color: COLORS.text,
            margin: 0,
            letterSpacing: "-0.04em",
          }}
        >
          Mochi CLI
        </h1>
        <p
          style={{
            fontSize: 24,
            color: COLORS.text,
            opacity: 0.7,
            marginTop: 10,
            fontWeight: 500,
          }}
        >
          Powerful automation for Mochi.cards
        </p>
      </div>

      {/* Terminal Window */}
      <div
        style={{
          position: "absolute",
          top: 220,
          left: "50%",
          transform: "translateX(-50%)",
          width: 1600,
        }}
      >
        <TerminalWindow>
          {/* Command Line */}
          <div style={{ marginBottom: 32, fontSize: 32 }}>
            <span style={{ color: COLORS.accent1, marginRight: 20, fontWeight: 700 }}>$</span>
            <TypewriterText
              text={currentSequence.command}
              startFrame={5}
              duration={typewriterDuration}
              frameInSequence={frameInSequence}
              color="#ffffff"
            />
            <Cursor
              visible={frameInSequence >= 5 && frameInSequence < typewriterDuration + 5}
              color={COLORS.accent1}
            />
          </div>

          {/* Output */}
          {frameInSequence >= outputDelay && (
            <OutputDisplay
              type={currentSequence.output as any}
              fadeInProgress={interpolate(
                frameInSequence,
                [outputDelay, outputDelay + outputFadeIn],
                [0, 1],
                { extrapolateLeft: "clamp", extrapolateRight: "clamp" }
              )}
            />
          )}
        </TerminalWindow>
      </div>

      {/* Feature label & Button */}
      <div
        style={{
          position: "absolute",
          bottom: 80,
          left: "50%",
          transform: "translateX(-50%)",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          gap: 30,
        }}
      >
        <div
          style={{
            fontSize: 28,
            color: COLORS.accent1,
            fontWeight: 700,
            letterSpacing: "0.05em",
            textTransform: "uppercase",
          }}
        >
          {currentSequence.description}
        </div>

        {/* Progress dots */}
        <div style={{ display: "flex", gap: 12 }}>
          {SEQUENCES.map((_, index) => {
            const isActive = index === currentSequenceIndex;
            return (
              <div
                key={index}
                style={{
                  width: isActive ? 40 : 12,
                  height: 12,
                  borderRadius: 6,
                  backgroundColor: isActive ? COLORS.accent1 : COLORS.text,
                  opacity: isActive ? 1 : 0.2,
                  transition: "all 0.3s ease",
                }}
              />
            );
          })}
        </div>
      </div>
    </AbsoluteFill>
  );
};

const Cursor: React.FC<{ visible: boolean; color: string }> = ({ visible, color }) => {
  const frame = useCurrentFrame();
  const blinkVisible = Math.floor(frame / 15) % 2 === 0;
  if (!visible) return null;
  return (
    <span
      style={{
        display: "inline-block",
        width: 4,
        height: 32,
        backgroundColor: color,
        marginLeft: 8,
        opacity: blinkVisible ? 1 : 0,
        verticalAlign: "middle",
      }}
    />
  );
};
