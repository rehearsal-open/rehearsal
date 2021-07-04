<!-- omit in toc -->
# REHEARSAL
```
 ______________________
[_/_\_/_\_/_\_/_\_/_\_/]
 \ /  /    _     \  \ /    ____        ____        __  _        ____        ____        ____        ____         ____        _
 \ / /    |o|     \ \ /   /____\      /____\      /_/  \\      /____\      /____\      /____\      /____\       /____\      //
  \ /   __/\\__    \ /   //____\\    //____      //_____\\    //____      //    \\    //____\\    //______     //    \\    //
  /_/      /\      \_\  // \ ___/   /______\  _ /________ \_ /______\  _ //______\\_ // \ ___/   /_______ \__ //______\\_ //
 //_\     /  \     /_\\ \\_ \ \____ \\_______// \\   ___/  / \\_______// \ ______  / \\_ \ \____ ________\  / \ ______  / \\_________
 //_\  __/\  /\__  /_\\  \_\ \____/  \_______/   \\  \____/   \_______/   \\ \____/   \_\ \____/ \_________/   \\ \____/   \________/

```

ハッカソン用の発表資料を公開しています。ぜひご覧ください。 [Google Slide]() [Youtube]() [niconico]()

## Contents


## Abstract/Target

### 動作テストを楽にしよう

　プログラムに限らず、何かプロダクトを開発するにあたり、一番時間をかけるべき場所はどこでしょうか？

　勿論、そのプロダクトが売りにしているところに時間をかけるべきではあります。しかしながら、消費者が見ているものは必ずしもその商品が売りにしているところとは限りません。要は全体を気にかけなければならないのです。とはいえ、一番売りにしたいところに時間を割くべきことには変わりありません。他は不具合が出ないような状態にしておけばよいのです。そのために、そのプロダクトが売りにしているところの次に割くべき時間は動作テストであるべきだと私は思います。

　確かに、世の中には様々なテストツールが存在します。しかし、それはどこかに特化した性能だといえるでしょう。それでもある程度のニーズには応えることができます。しかし、この世の中のニーズは刻一刻と変化しており、それに追随することもまた、大変です。そのニーズに追随するためにこのプロダクトを作りました。

> *よく勘違いされるのですが、このプログラムの思想はあくまでも「動作テストを制作しやすい・即時対応しやすい」ような環境を生み出すことに最大の意味があります。 **このプロダクトでは実際のテストコードを提供しないことにご注意ください。** その時代に応じて様々なテストコードへのアプローチがあったり、過去の遺産を流用したいというニーズに応えるものです。*

## Features

- 複数のプログラムを同時に実行し、それぞれの標準出力からデータを取得し、ほかのプログラムの標準入力にリレーさせます。
  - どのプログラムをどのような手段で実行し、標準出力をどのプログラムにリレーさせるかをまとめた `.yaml` ファイルを読み取ります。

### Future Features

- UART通信機能（組み込み開発テスト機能）
  - 個人的にロボコンに参加するのでそのセンサ値取得とグラフ化に使いたい（どんなグラフにするかはテスター部分で考えればよいので）
- アプリケーション間通信に使用しているアルゴリズムの全体的な刷新
  - 速度やメモリ使用量の観点から見直すべき点があまりに多すぎる
- 設定項目を充実させたい
  
### Try and Challenges

- `GoLang` の標準パッケージだけで極力開発しました。
- 開発設計にはクリーンアーキテクチャを採用し、効率よく追加機能を実装できるような状態にしました。

## Using Tools

![](icons/gopherbw.png) ![](icons/vscode.svg) ![](icons/github.svg)

- `GoLang` 
  - 標準パッケージ
  - [github.com/pkg/errors](https://github.com/pkg/errors)
  - [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3)
- `vscode` / `github`
  - いつものコードエディターです。

## Special Thanks
- デモ動画作成時に [chokudai](https://mobile.twitter.com/chokudai) 様より許可をいただきAtCoder内で公開されていた [AtCoder Hueristic Contest 003](https://atcoder.jp/contests/ahc003) とそのローカルテスタを使用させていただきました。
