#include <bits/stdc++.h>
using namespace std;

int main() {

    ios_base::sync_with_stdio(false);
    std::cin.tie(nullptr);

    long long N, pattern = 1, res;
    cin >> N;   // キーの長さを取得

    for (auto i = 0; i < N; i++ ){
        pattern *= 10;
    }

    // キー推測（愚直解）
    for (auto i = 0; i < pattern; i++ ) {

        char space[] = "0000000000";
        auto str = to_string(i);
        space[N - str.size()] = 0;

        cout << space << str << endl << flush;
        cin >> res;  // 結果を受け取る
        if (res == N) {
            return 0;
        }
    }
}

// int main() {
    
//     ios_base::sync_with_stdio(false);
//     std::cin.tie(nullptr);
//     for (auto i = 0; i < 10000; i++ ) {

//         char space[] = "0000000000";
//         auto str = to_string(i);
//         space[4 - str.size()] = 0;

//         cout << space << str << endl << flush;
//     }
// }
