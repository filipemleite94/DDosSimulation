#include<bits/stdc++.h>

using namespace std;

#define MS 1000000

int main(){
    map<long long,long long> m1,m2;
    unordered_map<string, long long> m;
    long long i,j,k,l;
    vector<long long> diff, lost;
    string s;
    int count(0);
    while(cin>>i){
        cin>>j;
        cin>>k;
        cin>>s;
        cin>>l;
        if(!m.count(s)) m[s] = 0;
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
    return 0;
}