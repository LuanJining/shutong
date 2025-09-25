import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'
import styleImport,{AntdResolve} from 'vite-plugin-style-import'

export default defineConfig({
  plugins: [
    react(),
    styleImport({
      resolve: [
        AntdResolve()
      ]
    })
  ],
  resolve:{
    alias:{
      "@":path.resolve(__dirname,'./src/')
    }
  },
//   server:{
//     host:'0.0.0.0'
//   }
})
