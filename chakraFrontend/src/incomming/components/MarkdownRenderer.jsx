import { marked } from 'marked';
import Prism from 'prismjs';
import { useEffect, useState } from 'react';

// Import the languages you need
import 'prismjs/components/prism-typescript';
import 'prismjs/components/prism-javascript';
import 'prismjs/components/prism-yaml';
import 'prismjs/components/prism-bash';

// Import a theme
import 'prismjs/themes/prism-dark.css';

const MarkdownRenderer = ({ content }) => {
  const [htmlContent, setHtmlContent] = useState('');

  useEffect(() => {
    if (content) {
      console.log('line 9 ', typeof content);

      // Configure marked with proper syntax highlighting
      marked.setOptions({
        breaks: true,
        gfm: true,
        highlight: (code, lang) => {
          if (lang && Prism.languages[lang]) {
            try {
              return Prism.highlight(code, Prism.languages[lang], lang);
            } catch (err) {
              console.warn('Prism highlighting failed:', err);
            }
          }
          // Fallback: escape HTML
          return code
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;');
        },
      });

      const html = marked(content);
      setHtmlContent(html);
    }
  }, [content]);

  useEffect(() => {
    // Re-highlight any code blocks that weren't caught by marked
    Prism.highlightAll();
  }, [htmlContent]);

  if (!content) {
    return <div className="text-gray-500 italic">No content available</div>;
  }

  return (
    <div
      className="markdown-content max-w-none text-gray-800 leading-relaxed"
      dangerouslySetInnerHTML={{ __html: htmlContent }}
    />
  );
};

export default MarkdownRenderer;
