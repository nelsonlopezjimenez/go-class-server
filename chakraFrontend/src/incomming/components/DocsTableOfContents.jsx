import { useEffect, useState } from 'react';

export default function DocsTableOfContents() {
  const [tocItems, setTocItems] = useState([]);

  useEffect(() => {
    // Extract headings from the current page
    const headings = document.querySelectorAll('h1, h2, h3, h4, h5, h6');
    const items = Array.from(headings).map((heading) => ({
      id:
        heading.id || heading.textContent.toLowerCase().replace(/[^\w]+/g, '-'),
      title: heading.textContent,
      level: Number.parseInt(heading.tagName.substring(1)),
    }));

    // Ensure headings have IDs
    headings.forEach((heading, index) => {
      if (!heading.id) {
        heading.id = items[index].id;
      }
    });

    setTocItems(items.filter((item) => item.level > 1)); // Skip h1
  }, []);

  return (
    <nav className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
      <h3 className="font-semibold mb-4 text-gray-900">Table of Contents</h3>
      {tocItems.length > 0 ? (
        <ul className="space-y-1">
          {tocItems.map((item, index) => (
            <li key={index}>
              <a
                href={`#${item.id}`}
                className="block text-sm py-1 px-2 text-gray-700 hover:text-indigo-600 transition-colors"
                style={{ marginLeft: `${(item.level - 2) * 12}px` }}
              >
                {item.title}
              </a>
            </li>
          ))}
        </ul>
      ) : (
        <p className="text-sm text-gray-500">No headings found</p>
      )}
    </nav>
  );
}
