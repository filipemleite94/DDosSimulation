#include<bits/stdc++.h>

using namespace std;

#define MS 1000000

int main(){
    map<long long,long long> m1,m2;
    unordered_map<string, long long> m;
    long long i,j,k,l;
    vector<long long> diff, lost, first;
    string s;
    int count(0);
    fstream fin("log");
    
    while(fin>>i){
        fin>>j;
        fin>>k;
        fin>>s;
        fin>>l;
        if(!m.count(s)){
            m[s] = 0;
            first.push_back(j/MS);
        }
        i/=MS*100;j/=MS*100;
        ++m1[i];
        ++m2[j];
        diff.push_back(k);
        if(count++%1000)    lost.back()+=l-m[s]-1;
        else lost.push_back(l-m[s]-1);
        m[s] = l;
    }
    long long min = m2.begin()->first;
    for(auto it:m1) cout<<(it.first-min)<<" "<<it.second<<endl;
    cout<<endl;
    for(auto it:m2) cout<<(it.first-min)<<" "<<it.second<<endl;
    cout<<endl;
    for(auto it:diff) cout<<it/1000<<endl;
    cout<<endl;
    for(auto it:lost) cout<<it<<endl;
    cout<<endl;
    min = first.front();
    for(auto it:first) cout<<it-min<<endl;
    return 0;
}