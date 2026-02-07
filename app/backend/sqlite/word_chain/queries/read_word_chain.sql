-- read_word_chain.sql
SELECT word_chain
FROM word_chain_puzzles
WHERE puzzle_date = :puzzle_date;
