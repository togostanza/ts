import Handlebars from 'handlebars/dist/handlebars';
import debounce from 'lodash.debounce';

function groupBy(array, func) {
  const ret = [];

  array.forEach((item) => {
    const key   = func(item);
    const entry = ret.filter((e) => e[0] === key)[0];

    if (entry) {
      entry[1].push(item);
    } else {
      ret.push([key, [item]]);
    }
  });

  return ret;
}

export default function initialize(descriptor) {
  return function Stanza(execute) {
    const development = descriptor.development;

    function template(name) {
      const t = descriptor.templates[name];
      if (!t) {
        throw new Error(`template "${name}" is not found`);
      }
      return t;
    }

    function createStanzaHelper(element) {
      const handlebars = Handlebars.create();

      return {
        root: element.shadowRoot,
        handlebars,

        query(params) {
          if (development) {
            console.log("query: called", params);
          }
          const t = template(params.template);
          const queryTemplate = handlebars.compile(t, {noEscape: true});
          const query = queryTemplate(params.parameters);
          const data = new URLSearchParams();
          data.set("query", query);

          if (development) {
            console.log("query: query built:\n" + query);
            console.log("query: sending to", params.endpoint);
          }

          // NOTE specifying Content-Type explicitly because some browsers sends `application/x-www-form-urlencoded;charset=UTF-8` without this, and some endpoints may not support this form.
          return fetch(params.endpoint, {
            method: params.method || "POST",
            headers: {
              "Content-Type": "application/x-www-form-urlencoded",
              "Accept": "application/sparql-results+json"
            },
            body: data,
          }).then((response) => {
            if (development) {
              console.log("query:", response.statusText, response);
            }

            return response.json();
          });
        },

        render(params) {
          if (development) {
            console.log("render: called", params)
          }

          const t = template(params.template);
          const htmlTemplate = handlebars.compile(t);
          const htmlFragment = htmlTemplate(params.parameters);

          if (development) {
            console.log("render: built:\n", htmlFragment)
          }

          const selector = params.selector || "main";
          element.shadowRoot.querySelector(selector).innerHTML = htmlFragment;

          if (development) {
            console.log("render: wrote to \"" + selector + "\"")
          }
        },

        select(selector) {
          return this.root.querySelector(selector);
        },

        selectAll(selector) {
          return this.root.querySelectorAll(selector);
        },

        grouping(rows, ...keys) {
          const normalizedKeys = keys.reduce((acc, key) => {
            if (key instanceof Array) {
              return acc.concat({key: key, alias: key.join('_')});
            } else if (key instanceof Object) {
              return acc.concat(key);
            } else {
              return acc.concat({key: key, alias: key});
            }
          }, []);

          return (function _grouping(rows, keys) {
            const [currentKey, ...remainKeys] = keys;

            function fetch(row, key) {
              return key instanceof Array ? key.map((k) => row[k]) : row[currentKey.key];
            }

            if (keys.length === 1) {
              return rows.map((row) => fetch(row, currentKey.key));
            }

            return groupBy(rows, (row) => {
              return fetch(row, currentKey.key);
            }).map(([currentValue, remainValues]) => {
              const nextKey = remainKeys[0];

              return {
                [currentKey.alias]: currentValue,
                [nextKey.alias]:    _grouping(remainValues, remainKeys)
              };
            });
          })(rows, normalizedKeys);
        },

        groupBy,

        unwrapValueFromBinding(queryResult) {
          const bindings = queryResult.results.bindings;

          return bindings.map((binding) => {
            const ret = {};

            Object.keys(binding).forEach((key) => {
              ret[key] = binding[key].value;
            });

            return ret;
          });
        }
      };
    }

    const update = debounce((element) => {
      const params = descriptor.parameters.reduce((acc, key) => Object.assign(acc, {[key]: element.getAttribute(key)}), {});

      execute(createStanzaHelper(element), params);
    }, 50);

    class StanzaElement extends HTMLElement {
      constructor() {
        super();

        const shadow = this.attachShadow({mode: "open"});
        const main = document.createElement("main");
        shadow.appendChild(main);

        update(this);
      }

      static get observedAttributes() {
        return descriptor.parameters;
      }

      attributeChangedCallback(attrName, oldVal, newVal) {
        if (!descriptor.parameters.includes(attrName)) { return; }

        update(this);
      }
    }

    if ('customElements' in window && !window.customElements.get(descriptor.elementName)) {
      window.customElements.define(descriptor.elementName, StanzaElement);
    }
  }
}
