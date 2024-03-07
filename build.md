Para criar pacotes `.deb` (para distribuições baseadas no Debian, como Ubuntu) e `.pacman` (para distribuições baseadas no Arch Linux, como Manjaro), você pode usar ferramentas específicas para cada uma delas.

### Para criar um pacote `.deb` (Debian/Ubuntu):

1. **Instale o pacote `build-essential`:**
   ```bash
   sudo apt-get update
   sudo apt-get install build-essential
   ```

2. **Compile o seu programa Go e crie um executável:**
   ```bash
   go build -o meu_programa
   ```

3. **Crie a estrutura do diretório do pacote:**
   ```bash
   mkdir -p ~/meu_programa/DEBIAN
   ```

4. **Crie um arquivo `control` dentro do diretório `DEBIAN`:**
   ```plaintext
   Package: meu_programa
   Version: 1.0
   Architecture: amd64
   Maintainer: Seu Nome <seu@email.com>
   Description: Sua descrição aqui.
   ```

5. **Mova o executável para o diretório `usr/bin`:**
   ```bash
   mkdir -p ~/meu_programa/usr/bin
   mv meu_programa ~/meu_programa/usr/bin/
   ```

6. **Defina as permissões corretas:**
   ```bash
   chmod +x ~/meu_programa/DEBIAN/control
   chmod +x ~/meu_programa/usr/bin/meu_programa
   ```

7. **Crie o pacote `.deb`:**
   ```bash
   dpkg-deb --build ~/meu_programa
   ```

   Isso criará um arquivo `.deb` no diretório atual.

### Para criar um pacote `.pacman` (Arch Linux/Manjaro):

1. **Instale o pacote `base-devel`:**
   ```bash
   sudo pacman -S base-devel
   ```

2. **Compile o seu programa Go e crie um executável:**
   ```bash
   go build -o meu_programa
   ```

3. **Crie uma estrutura de diretório do PKGBUILD:**
   ```bash
   mkdir -p ~/meu_programa/PKGBUILD
   ```

4. **Crie um arquivo `PKGBUILD` dentro do diretório `meu_programa`:**
   ```bash
   nano ~/meu_programa/PKGBUILD
   ```
   ```bash
   pkgname=meu_programa
   pkgver=1.0
   pkgrel=1
   arch=('x86_64')
   url='https://seusite.com'
   license=('MIT')
   depends=('go')
   source=("$pkgname-$pkgver.tar.gz::$url/archive/v$pkgver.tar.gz")
   sha256sums=('SKIP')

   build() {
       cd "$srcdir/$pkgname-$pkgver"
       go build -o $pkgname
   }

   package() {
       cd "$srcdir/$pkgname-$pkgver"
       install -Dm755 $pkgname "$pkgdir/usr/bin/$pkgname"
   }
   ```

   Modifique as informações conforme necessário.

5. **Compacte o código-fonte:**
   ```bash
   tar -czvf ~/meu_programa-$pkgver.tar.gz meu_programa
   ```

6. **Crie o pacote `.pacman`:**
   ```bash
   makepkg -si
   ```

   Isso instalará o pacote ou você pode usar o pacote gerado diretamente.

Lembre-se de que esses são procedimentos básicos e você deve ajustá-los conforme necessário para o seu caso específico. Certifique-se de incluir todas as dependências necessárias para que o seu programa funcione corretamente nas diferentes distribuições.