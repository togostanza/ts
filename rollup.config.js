import babel from 'rollup-plugin-babel';
import commonjs from 'rollup-plugin-commonjs';
import resolve from 'rollup-plugin-node-resolve';
import { uglify } from 'rollup-plugin-uglify';

export default {
  input: 'provider/assets-src/js/stanza.js',
  output: {
    file: 'provider/assets/js/stanza.js',
    format: 'iife',
    name: 'Stanza',
    sourcemap: true
  },
  plugins: [
    babel({
      exclude: 'node_modules/**'
    }),
    resolve(),
    commonjs({
      include: 'node_modules/**'
    }),
    uglify()
  ]
};
