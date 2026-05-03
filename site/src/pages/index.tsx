import React, { StrictMode } from 'react';
import { Installation } from '@/src/components/Installation';
import { Hero } from '@/src/components/Hero';
import Head from '../components/Head';
import { Footer } from '../components/Footer';
import { CodeBlock } from '@/src/components/CodeBlock';

export default function Home() {
  return (
    <div className="wrapper dark">
      <Head />
      <main className="container">
        <Hero />
        <div className="content">
          <Installation />
          <div>
            <h2>Quick Start</h2>
            <p>Stage, generate, commit, and push — in one command.</p>
            <CodeBlock initialHeight={120}>{`gli commit
# AI generates a conventional commit message
# You review → commit & push`}
            </CodeBlock>
          </div>
          <div>
            <h2>Commands</h2>
            <p>Five commands. No fluff.</p>
            <CodeBlock initialHeight={180}>{`gli commit              # AI-powered commit + push
gli branch -c <name>   # Create & push branch
gli branch -s <name>   # Switch branch
gli me                  # Your GitHub profile
gli profile <user>      # Any GitHub profile`}
            </CodeBlock>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  );
}
