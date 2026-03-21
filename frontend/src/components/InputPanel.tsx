import { useRef, useState } from "react";

type Props = {
  rawInput: string;
  setRawInput: (v: string) => void;
  patientName: string;
  setPatientName: (v: string) => void;
  patientAge: number;
  setPatientAge: (v: number) => void;
  onProcess: () => void;
  processing: boolean;
};

type SpeechRecognitionConstructor = new () => {
  continuous: boolean;
  interimResults: boolean;
  lang: string;
  onresult: ((event: { results: ArrayLike<ArrayLike<{ transcript: string }>> }) => void) | null;
  onerror: (() => void) | null;
  start: () => void;
  stop: () => void;
};

export default function InputPanel(props: Props) {
  const [listening, setListening] = useState(false);
  const recognitionRef = useRef<InstanceType<SpeechRecognitionConstructor> | null>(null);

  function startVoiceInput() {
    const SpeechRecognition =
      (window as unknown as { SpeechRecognition?: SpeechRecognitionConstructor; webkitSpeechRecognition?: SpeechRecognitionConstructor }).SpeechRecognition ||
      (window as unknown as { SpeechRecognition?: SpeechRecognitionConstructor; webkitSpeechRecognition?: SpeechRecognitionConstructor }).webkitSpeechRecognition;

    if (!SpeechRecognition) {
      alert("Speech recognition is not supported in this browser.");
      return;
    }

    if (!recognitionRef.current) {
      const recognition = new SpeechRecognition();
      recognition.continuous = true;
      recognition.interimResults = true;
      recognition.lang = "en-US";
      recognition.onresult = (event) => {
        const combined = Array.from(event.results)
          .map((res) => res[0].transcript)
          .join(" ");
        props.setRawInput(combined.trim());
      };
      recognition.onerror = () => setListening(false);
      recognitionRef.current = recognition;
    }

    recognitionRef.current.start();
    setListening(true);
  }

  function stopVoiceInput() {
    recognitionRef.current?.stop();
    setListening(false);
  }

  return (
    <div className="panel">
      <h2>Unified Clinical Input</h2>
      <label>
        Patient Name
        <input value={props.patientName} onChange={(e) => props.setPatientName(e.target.value)} placeholder="John Doe" />
      </label>
      <label>
        Patient Age
        <input type="number" value={props.patientAge} onChange={(e) => props.setPatientAge(Number(e.target.value))} />
      </label>
      <label>
        Clinical Notes (text or voice)
        <textarea
          rows={10}
          value={props.rawInput}
          onChange={(e) => props.setRawInput(e.target.value)}
          placeholder="Type or dictate clinical observations, medications, and tests..."
        />
      </label>
      <div style={{ display: "flex", gap: 8, marginBottom: 10 }}>
        {!listening ? (
          <button type="button" onClick={startVoiceInput}>Start Voice</button>
        ) : (
          <button type="button" onClick={stopVoiceInput}>Stop Voice</button>
        )}
        <button className="process-btn" type="button" onClick={props.onProcess} disabled={props.processing}>
          {props.processing ? "Processing..." : "Process with AI"}
        </button>
      </div>
      <p>
        Save is disabled until AI output is reviewed and edited in the verification panel.
      </p>
    </div>
  );
}
